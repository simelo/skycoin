typedef struct{
} visor__unconfirmedTxns;
typedef struct{
} visor__txUnspents;
typedef struct{
    coin__Transaction Txn;
    GoInt64_ Received;
    GoInt64_ Checked;
    GoInt64_ Announced;
    GoInt8_ IsValid;
} visor__UnconfirmedTxn;