#include <stdlib.h>
#include <stdio.h>
#include <sqlite3.h>
#include <string.h>

sqlite3* prepareReading(char * filename, int n)
{
   sqlite3* DB;
   int exit = 0;
   char* messaggeError;
   sqlite3_stmt *stm;

   sqlite3_config(SQLITE_CONFIG_URI,1);
   exit = sqlite3_open(filename, &DB);
   if (exit != SQLITE_OK) {
       printf("failed to open!!\n");
       sqlite3_close(DB);
       return NULL;
   }
   char *sql = "create table t(i int);"
	               "begin;";
   exit = sqlite3_exec(DB, sql, NULL, 0, &messaggeError);
   if (exit != SQLITE_OK) {
       printf("failed to create table!!: %s\n", messaggeError);
       sqlite3_close(DB);
       return NULL;
   }
   sql = "insert into t values(?)";
   if (sqlite3_prepare_v2(DB, sql, -1, &stm, NULL) != SQLITE_OK) {
       printf("failed to prepare inserting!!\n");
       sqlite3_close(DB);
       return NULL;
   }
   int i;
   for (i = 0; i < n; i++){
       if (sqlite3_bind_int(stm, 1, i) != SQLITE_OK) {
           printf("failed to bind value!!\n");
           sqlite3_close(DB);
           return NULL;
       }
       if (sqlite3_step(stm) != SQLITE_DONE) {
           printf("failed to create table!!\n");
           sqlite3_close(DB);
           return NULL;
       }
       sqlite3_reset(stm);
   }
   sqlite3_finalize(stm);
   sql = "commit";
   exit = sqlite3_exec(DB, sql, NULL, 0, &messaggeError);
   if (exit != SQLITE_OK) {
       printf("failed to commit!!: %s\n", messaggeError);
       sqlite3_close(DB);
       return NULL;
   }
   return DB;
}
void reading(sqlite3 *DB, int n)
{
   int exit = 0;
   sqlite3_stmt *stm;

   char *sql = "select * from t";
   if (sqlite3_prepare_v2(DB, sql, -1, &stm, NULL) != SQLITE_OK) {
       printf("failed prepare select!!2:%s\n", sqlite3_errmsg(DB));
       sqlite3_close(DB);
       return;
   }
   int i;
   for (i = 0; i < n; i++){
       if (sqlite3_step(stm) != SQLITE_ROW) {
           printf("failed to read from table!!:%s\n", sqlite3_errmsg(DB));
           sqlite3_close(DB);
           return;
       }
       int res = sqlite3_column_int(stm, 0);
   }
   sqlite3_finalize(stm);
   return;
}


sqlite3* prepareInsertComparative(char * filename, int n)
{
   sqlite3* DB;
   int exit = 0;
   char* messaggeError;

   sqlite3_config(SQLITE_CONFIG_URI,1);
   exit = sqlite3_open(filename, &DB);
   if (exit != SQLITE_OK) {
       printf("failed to open!!\n");
       sqlite3_close(DB);
       return NULL;
   }
   char *sql = "create table t(i int);";
   exit = sqlite3_exec(DB, sql, NULL, 0, &messaggeError);
   if (exit != SQLITE_OK) {
       printf("failed to create table!!: %s\n", messaggeError);
       sqlite3_close(DB);
       return NULL;
   }
    return DB;
}
void insertComparative(sqlite3 *DB, int n)
{

   int exit = 0;
   char* messaggeError;
   sqlite3_stmt *stm;
   char *sql = "begin;"
	               "delete from t;";
   exit = sqlite3_exec(DB, sql, NULL, 0, &messaggeError);
   if (exit != SQLITE_OK) {
       printf("start transaction!!: %s\n", messaggeError);
       sqlite3_close(DB);
       return;
   }

   sql = "insert into t values(?)";
   if (sqlite3_prepare_v2(DB, sql, -1, &stm, NULL) != SQLITE_OK) {
       printf("failed to prepare inserting!!\n");
       sqlite3_close(DB);
       return;
   }
   int i;
   for (i = 0; i < n; i++){
       if (sqlite3_bind_int(stm, 1, i) != SQLITE_OK) {
           printf("failed to bind value!!\n");
           sqlite3_close(DB);
           return;
       }
       if (sqlite3_step(stm) != SQLITE_DONE) {
           printf("failed to create table!!\n");
           sqlite3_close(DB);
           return;
       }
       sqlite3_reset(stm);
   }
   sqlite3_finalize(stm);
   sql = "commit";
   exit = sqlite3_exec(DB, sql, NULL, 0, &messaggeError);
   if (exit != SQLITE_OK) {
       printf("failed to commit!!: %s\n", messaggeError);
       sqlite3_close(DB);
       return;
   }
   return;
} 