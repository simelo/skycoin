
#include <stdio.h>

#include "cipher.testsuite.testsuite.go.h"

TestSuite(cipher_testsuite, .init = setup, .fini = teardown);

Test(cipher_testsuite, TestManyAddresses) {
  SeedTestDataJSON dataJSON;
  SeedTestData data;
  GoUint32 err;

  json_value* json = loadGoldenFile(MANY_ADDRESSES_FILENAME);
  cr_assert(json != NULL, "Error loading file");
  jsonToSeedTestData(json, &dataJSON);
  err = SeedTestDataFromJSON(&dataJSON, &data);
  cr_assert(err == SKY_OK, "Deserializing seed test data from JSON ... %d", err);
  ValidateSeedData(&data, NULL);
}
