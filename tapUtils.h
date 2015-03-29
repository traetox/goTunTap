#ifndef _TAPUTILS_H_
#define _TAPUTILS_H_

int StartTap(char *name);
int StopTap(int sock, char* name);
int AddTapToBridge(char* bridge, char* tap);
int RemoveTapFromBridge(char* bridge, char* tap);
int CreateBridge(char* bridge);
int DeleteBridge(char* bridge);
int CheckBridge(char* bridge);

#endif
