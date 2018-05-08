#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#include <criterion/criterion.h>
#include <criterion/new/assert.h>

#include "libskycoin.h"
#include "skyerrors.h"
#include "skystring.h"
#include "skytest.h"

#define BUFFER_SIZE 1024
#define PASSWORD1 "pwd"
#define PASSWORD2 "key"
#define PASSWORD3 "9JMkCPphe73NQvGhmab"
#define SHA256XORDATALENGTHSIZE 4
#define SHA256XORBLOCKSIZE 32
#define SHA256XORCHECKSUMSIZE 32
#define SHA256XORNONCESIZE 32


typedef struct{
	int 		dataLength;
	GoSlice* 	pwd;
	int			success;
} TEST_DATA;

/*
// PutUvarint encodes a uint64 into buf and returns the number of bytes written.
// If the buffer is too small, PutUvarint will panic.
func PutUvarint(buf []byte, x uint64) int {
	i := 0
	for x >= 0x80 {
		buf[i] = byte(x) | 0x80
		x >>= 7
		i++
	}
	buf[i] = byte(x)
	return i + 1
}

// PutVarint encodes an int64 into buf and returns the number of bytes written.
// If the buffer is too small, PutVarint will panic.
func PutVarint(buf []byte, x int64) int {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	return PutUvarint(buf, ux)
}

func hashKeyIndexNonce(key []byte, index int64, nonceHash cipher.SHA256) cipher.SHA256 {
	// convert index to 256bit number
	indexBytes := make([]byte, 32)
	binary.PutVarint(indexBytes, index)

	// hash(index, nonceHash)
	indexNonceHash := cipher.SumSHA256(append(indexBytes, nonceHash[:]...))

	// hash(hash(password), indexNonceHash)
	var keyHash cipher.SHA256
	copy(keyHash[:], key[:])
	return cipher.AddSHA256(keyHash, indexNonceHash)
}
{
		FILE *fp;
		fp = fopen("test.txt", "w+");
		for( int i = 0; i < 32; i++){
			if ( buff[i] > 0 )
				fprintf(fp, "%d\n", buff[i]);
		}
		fclose(fp);
	}
*/

int putUvarint(GoSlice* buf , GoUint64 x){
	int i = 0;
	while( x >= 0x80 && i < buf->cap) {
		((unsigned char*)buf->data)[i] = ((GoUint8)x) | 0x80;
		x >>= 7;
		i++;
	}
	if( i < buf->cap ){
		((unsigned char*)buf->data)[i] = (GoUint8)(x);
		buf->len = i + 1;
	} else {
		buf->len = i;
	}
	return buf->len;
}

int putVarint(GoSlice* buf , GoInt64 x){
	GoUint64 ux = (GoUint64)x << 1;
	if ( x < 0 ) {
		ux = ~ux;
	}
	return putUvarint(buf, ux);
}

void hashKeyIndexNonce(GoSlice_ key, GoInt64 index, 
	cipher__SHA256 *nonceHash, cipher__SHA256 *resultHash){
	GoUint32 errcode;
	int length = 32 + sizeof(cipher__SHA256);
	unsigned char buff[length];
	GoSlice slice = {buff, 0, length};
	memset(buff, 0, length * sizeof(char));
	putVarint( &slice, index );
	memcpy(buff + 32, *nonceHash, sizeof(cipher__SHA256));
	slice.len = length;
	cipher__SHA256 indexNonceHash;
	errcode = SKY_cipher_SumSHA256(slice, indexNonceHash);
	cr_assert(errcode == SKY_OK, "SKY_cipher_SumSHA256 failed. Error calculating hash");
	SKY_cipher_AddSHA256(key.data, &indexNonceHash, resultHash);
	cr_assert(errcode == SKY_OK, "SKY_cipher_AddSHA256 failed. Error adding hashes");
}

void makeEncryptedData(GoSlice data, GoUint32 dataLength, GoSlice pwd, GoSlice* encrypted){
	GoUint32 fullLength = dataLength + SHA256XORDATALENGTHSIZE;
	GoUint32 n = fullLength / SHA256XORBLOCKSIZE;
	GoUint32 m = fullLength % SHA256XORBLOCKSIZE;
	GoUint32 errcode;
	
	if( m > 0 ){
		fullLength += SHA256XORBLOCKSIZE - m;
	}
	cr_assert(SHA256XORBLOCKSIZE == sizeof(cipher__SHA256), "Size of SHA256 block size different that cipher.SHA256 struct");
	fullLength += SHA256XORBLOCKSIZE;
	char* buffer = malloc(fullLength);
	cr_assert(buffer != NULL, "Couldn\'t allocate buffer");
	//Add data length to the beginning, saving space for the checksum
	for(int i = 0; i < SHA256XORDATALENGTHSIZE; i++){
		int shift = i * 8;
		buffer[i + SHA256XORBLOCKSIZE] = (dataLength & (0xFF << shift)) >> shift;
	}
	//Add the data
	memcpy(buffer + SHA256XORDATALENGTHSIZE + SHA256XORBLOCKSIZE, 
		data.data, dataLength);
	/*for(int i = 0; i < dataLength; i++){
		buffer[i + SHA256XORDATALENGTHSIZE + SHA256XORBLOCKSIZE] = ((char*)data.data)[i];
	}*/
	//Add padding
	for(int i = dataLength + SHA256XORDATALENGTHSIZE + SHA256XORBLOCKSIZE; i < fullLength; i++){
		buffer[i] = 0;
	}
	//Buffer with space for the checksum, then data length, then data, and then padding
	GoSlice _data = {buffer + SHA256XORBLOCKSIZE, 
		fullLength - SHA256XORBLOCKSIZE, 
		fullLength - SHA256XORBLOCKSIZE};
	//GoSlice _hash = {buffer, 0, SHA256XORBLOCKSIZE};
	errcode = SKY_cipher_SumSHA256(_data, buffer);
	cr_assert(errcode == SKY_OK, "SKY_cipher_SumSHA256 failed. Error calculating hash");
	char bufferNonce[SHA256XORNONCESIZE];
	GoSlice sliceNonce = {bufferNonce, 0, SHA256XORNONCESIZE};
	randBytes(&sliceNonce, SHA256XORNONCESIZE);
	cipher__SHA256 hashNonce;
	errcode = SKY_cipher_SumSHA256(sliceNonce, hashNonce);
	cr_assert(errcode == SKY_OK, "SKY_cipher_SumSHA256 failed. Error calculating hash for nonce");
	char bufferHash[BUFFER_SIZE];
	coin__UxArray hashPassword = {bufferHash, 0, BUFFER_SIZE};
	errcode = SKY_secp256k1_Secp256k1Hash(pwd, &hashPassword);
	cr_assert(errcode == SKY_OK, "SKY_secp256k1_Secp256k1Hash failed. Error calculating hash for password");
	cipher__SHA256 h;
	
	
	int fullDestLength = fullLength + sizeof(cipher__SHA256) + SHA256XORNONCESIZE;
	int destBufferStart = sizeof(cipher__SHA256) + SHA256XORNONCESIZE;
	char* dest_buffer = malloc(fullDestLength);
	cr_assert(dest_buffer != NULL, "Couldn\'t allocate result buffer");
	for(int i = 0; i < n; i++){
		hashKeyIndexNonce(hashPassword, i, &hashNonce, &h);
		cipher__SHA256* pBuffer = (cipher__SHA256*)(buffer + i *SHA256XORBLOCKSIZE);
		cipher__SHA256* xorResult = (cipher__SHA256*)(dest_buffer + destBufferStart + i *SHA256XORBLOCKSIZE);
		SKY_cipher_SHA256_Xor(pBuffer, &h, xorResult);
	}
	// Prefix the nonce
	memcpy(dest_buffer + sizeof(cipher__SHA256), bufferNonce, SHA256XORNONCESIZE);
	// Calculates the checksum
	GoSlice nonceAndDataBytes = {dest_buffer + sizeof(cipher__SHA256), 
								fullLength + SHA256XORNONCESIZE, 
								fullLength + SHA256XORNONCESIZE
						};
	cipher__SHA256* checksum = (cipher__SHA256*)dest_buffer;
	errcode = SKY_cipher_SumSHA256(nonceAndDataBytes, checksum);
	cr_assert(errcode == SKY_OK, "SKY_cipher_SumSHA256 failed. Error calculating final checksum");
}

Test(cipher_encrypt_sha256xor, TestSha256XorEncrypt){
	unsigned char buff[BUFFER_SIZE];
	unsigned char encryptedBuffer[BUFFER_SIZE];
	unsigned char encryptedText[BUFFER_SIZE];
	GoSlice data = { buff, 0, BUFFER_SIZE };
	GoSlice encrypted = { encryptedBuffer, 0, BUFFER_SIZE };
	GoSlice pwd1 = { PASSWORD1, strlen(PASSWORD1), strlen(PASSWORD1) };
	GoSlice pwd2 = { PASSWORD2, strlen(PASSWORD2), strlen(PASSWORD2) };
	GoSlice pwd3 = { PASSWORD3, strlen(PASSWORD3), strlen(PASSWORD3) };
	GoSlice nullPwd = {NULL, 0, 0};
	GoUint32 errcode;
	
	TEST_DATA test_data[] = {
		{1, &nullPwd, 0},
		{1, &pwd2, 1},
		{2, &pwd1, 1},
		{32, &pwd1, 1},
		{64, &pwd3, 1},
		{65, &pwd3, 1},
	};
	
	encrypt__Sha256Xor encryptSettings = {};
	
	for(int i = 0; i < 6; i++){
		randBytes(&data, test_data[i].dataLength);
		errcode = SKY_encrypt_Sha256Xor_Encrypt(&encryptSettings, data, *(test_data[i].pwd), &encrypted);
		if( test_data[i].success ){
			cr_assert(errcode == SKY_OK, "SKY_encrypt_Sha256Xor_Encrypt failed.");
		} else {
			cr_assert(errcode != SKY_OK, "SKY_encrypt_Sha256Xor_Encrypt with null pwd.");
		}
		if( errcode == SKY_OK ){
			cr_assert(encrypted.cap > 0, "Buffer for encrypted data is too short");
			cr_assert(encrypted.len < BUFFER_SIZE, "Too large encrypted data");
			((char*)encrypted.data)[encrypted.len] = 0;
			
			int n = (SHA256XORDATALENGTHSIZE + test_data[i].dataLength) / SHA256XORBLOCKSIZE;
			int m = (SHA256XORDATALENGTHSIZE + test_data[i].dataLength) % SHA256XORBLOCKSIZE;
			if ( m > 0 ) {
				n++;
			}
			
			int real_size;
			base64_decode_binary((const unsigned char*)encrypted.data, 
				encrypted.len, encryptedText, &real_size, BUFFER_SIZE);
			int totalEncryptedDataLen = SHA256XORCHECKSUMSIZE + SHA256XORNONCESIZE + 32 + n*SHA256XORBLOCKSIZE; // 32 is the hash data length
			
			cr_assert(totalEncryptedDataLen == real_size, "SKY_encrypt_Sha256Xor_Encrypt failed, encrypted data length incorrect.");
			cr_assert(SHA256XORCHECKSUMSIZE == sizeof(cipher__SHA256), "Size of SHA256 struct different than size in constant declaration");
			cipher__SHA256 enc_hash;
			cipher__SHA256 cal_hash;
			for(int j = 0; j < SHA256XORCHECKSUMSIZE; j++){
				enc_hash[j] = (GoUint8_)encryptedText[j];
			}
			int len_minus_checksum = real_size - SHA256XORCHECKSUMSIZE;
			GoSlice slice = {&encryptedText[SHA256XORCHECKSUMSIZE], len_minus_checksum, len_minus_checksum};
			SKY_cipher_SumSHA256(slice, &cal_hash);
			int equal = 1;
			for(int j = 0; j < SHA256XORCHECKSUMSIZE; j++){
				if(enc_hash[j] != cal_hash[j]){
					equal = 0;
					break;
				}
			}
			cr_assert(equal == 1, "SKY_encrypt_Sha256Xor_Encrypt failed, incorrect hash sum.");
		}
	}
	
	for(int i = 33; i <= 64; i++){
		randBytes(&data, i);
		errcode = SKY_encrypt_Sha256Xor_Encrypt(&encryptSettings, data, pwd1, &encrypted);
		cr_assert(errcode == SKY_OK, "SKY_encrypt_Sha256Xor_Encrypt failed.");
		cr_assert(encrypted.cap > 0, "Buffer for encrypted data is too short");
		cr_assert(encrypted.len < BUFFER_SIZE, "Too large encrypted data");
		((char*)encrypted.data)[encrypted.len] = 0;
		
		int n = (SHA256XORDATALENGTHSIZE + i) / SHA256XORBLOCKSIZE;
		int m = (SHA256XORDATALENGTHSIZE + i) % SHA256XORBLOCKSIZE;
		if ( m > 0 ) {
			n++;
		}
		
		int real_size;
		base64_decode_binary((const unsigned char*)encrypted.data, 
			encrypted.len, encryptedText, &real_size, BUFFER_SIZE);
		int totalEncryptedDataLen = SHA256XORCHECKSUMSIZE + SHA256XORNONCESIZE + 32 + n*SHA256XORBLOCKSIZE; // 32 is the hash data length
		
		cr_assert(totalEncryptedDataLen == real_size, "SKY_encrypt_Sha256Xor_Encrypt failed, encrypted data length incorrect.");
		cr_assert(SHA256XORCHECKSUMSIZE == sizeof(cipher__SHA256), "Size of SHA256 struct different than size in constant declaration");
		cipher__SHA256 enc_hash;
		cipher__SHA256 cal_hash;
		for(int j = 0; j < SHA256XORCHECKSUMSIZE; j++){
			enc_hash[j] = (GoUint8_)encryptedText[j];
		}
		int len_minus_checksum = real_size - SHA256XORCHECKSUMSIZE;
		GoSlice slice = {&encryptedText[SHA256XORCHECKSUMSIZE], len_minus_checksum, len_minus_checksum};
		SKY_cipher_SumSHA256(slice, cal_hash);
		int equal = 1;
		for(int j = 0; j < SHA256XORCHECKSUMSIZE; j++){
			if(enc_hash[j] != cal_hash[j]){
				equal = 0;
				break;
			}
		}
		cr_assert(equal == 1, "SKY_encrypt_Sha256Xor_Encrypt failed, incorrect hash sum.");
		
	}
}

Test(cipher_encrypt_sha256xor, TestSha256XorDecrypt){
	unsigned char buff[BUFFER_SIZE];
	GoSlice data = {buff, 0, BUFFER_SIZE};
	GoSlice pwd = { PASSWORD1, strlen(PASSWORD1), strlen(PASSWORD1) };
	randBytes(&data, 32);
	GoSlice encrypted;
	makeEncryptedData(data, 32, pwd, &encrypted);
}

