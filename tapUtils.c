#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<sys/types.h>
#include<sys/stat.h>
#include<unistd.h>
#include<fcntl.h>
#include<sys/ioctl.h>
#include<sys/fcntl.h>
#include<linux/if_tun.h>
#include<net/if.h>
#include<sys/socket.h>

/* some manual defines because our libc is so damn old */
#ifndef SIOCBRADDBR
#define SIOCBRADDBR 0x89a0
#endif
#ifndef SIOCBRDELBR
#define SIOCBRDELBR 0x89a1
#endif
#ifndef SIOCBRADDIF
#define SIOCBRADDIF 0x89a2
#endif
#ifndef SIOCBRDELIF
#define SIOCBRDELIF 0x89a3
#endif

int setClearIfReqFlags(char* dev, int flags, int set) {
	struct ifreq ifr;
	int sock;
	memset(&ifr, 0, sizeof(ifr));
	strncpy(ifr.ifr_name, dev, sizeof(ifr.ifr_name));
	sock = socket(PF_INET, SOCK_DGRAM, 0);
	if(sock < 0) {
		return -1;
	}

	if(ioctl(sock, SIOCGIFFLAGS, &ifr) < 0) {
		close(sock);
		return -1;
	}

	if(set) {
		ifr.ifr_flags |= flags;
	} else {
		ifr.ifr_flags &= (~flags);
	}
	if(ioctl(sock, SIOCSIFFLAGS, &ifr) < 0) {
		close(sock);
		return -1;
	}
	close(sock);
	return 0;
}

int ifup(char *dev) {
	int flags = (IFF_UP|IFF_BROADCAST|IFF_RUNNING|IFF_MULTICAST);
	return setClearIfReqFlags(dev, flags, 1);
}

int ifdown(char *dev) {
	int flags = (IFF_UP|IFF_RUNNING);
	return setClearIfReqFlags(dev, flags, 0);
}

int tun_alloc(char *dev) {
      struct ifreq ifr;
      int fd, err;

      if( (fd = open("/dev/net/tun", O_RDWR)) < 0 )
         return -1;

      memset(&ifr, 0, sizeof(ifr));

      /* Flags: IFF_TUN   - TUN device (no Ethernet headers) 
       *        IFF_TAP   - TAP device  
       *
       *        IFF_NO_PI - Do not provide packet information  
       */ 
      ifr.ifr_flags = IFF_TAP; 
      if( *dev )
         strncpy(ifr.ifr_name, dev, IFNAMSIZ);

      if( (err = ioctl(fd, TUNSETIFF, (void *) &ifr)) < 0 ){
         close(fd);
         return err;
      }
      return fd;
}

int tun_dealloc(int tapfd) {
	int r = ioctl(tapfd, TUNSETPERSIST, 0);
	close(tapfd);
	return r;
}

int StartTap(char *name) {
	int tapfd;
	if((tapfd = tun_alloc(name)) <= 0) {
		printf("Failed to allocate tap %s\n", name);
		return -1;
	}
	if(ifup(name) != 0) {
		printf("Failed to bring up interface %s\n", name);
		return -1;
	}
	return tapfd;
}

int StopTap(int sock, char* name) {
	if(ifdown(name) != 0) {
		printf("Failed to bring down interface %s\n", name);
		return -1;
	}
	if(tun_dealloc(sock) != 0) {
		printf("Failed to dealloc tap %s\n", name);
		return -1;
	}
	return 0;
}

int CheckBridge(char* bridge) {
	char buff[256];
	struct stat s;
	snprintf(buff, sizeof(buff), "/sys/class/net/%s/bridge", bridge);
	
	if(stat(buff, &s) != 0) {
		return -1;
	}
	if(!S_ISDIR(s.st_mode)) {
		return -1;
	} else {
		return ifup(bridge);
	}
}

int CreateBridge(char* bridge) {
	int sock = socket(AF_LOCAL, SOCK_STREAM, 0);
	int br = -1;
	if(sock <= 0) {
		return -1;
	}
	br = ioctl(sock, SIOCBRADDBR, bridge);
	close(sock);
	if(br < 0 ) {
		return -1;
	}
	return ifup(bridge);
}

int DeleteBridge(char* bridge) {
	int sock = -1;
	int br = -1;
	/* if it doesn't exist then its deleted! */
	if(CheckBridge(bridge) < 0) {
		return 0;
	}

	sock = socket(AF_LOCAL, SOCK_STREAM, 0);
	if(sock <= 0) {
		return -1;
	}
	if(ifdown(bridge) != 0) {
		return -1;
	}
	br = ioctl(sock, SIOCBRDELBR, bridge);
	close(sock);
	return br < 0 ? -1 : 0;
}

int RemoveTapFromBridge(char* bridge, char* tap) {
	struct ifreq ir;
	int err;
	int sock;
	int ifindex = if_nametoindex(tap);
	if(ifindex == 0) {
		return -1;
	}

	if(CheckBridge(bridge) != 0) {
		return -1;
	}
	
	sock = socket(AF_LOCAL, SOCK_STREAM, 0);
	if(sock <= 0) {
		return -1;
	}

	strncpy(ir.ifr_name, bridge, IFNAMSIZ);
	ir.ifr_ifindex = ifindex;
	if(ioctl(sock, SIOCBRDELIF, &ir) < 0) {
		close(sock);
		return -1;
	}
	close(sock);
	return 0;
}

int AddTapToBridge(char* bridge, char* tap) {
	struct ifreq ir;
	int err;
	int sock;
	int ifindex = if_nametoindex(tap);
	if(ifindex == 0) {
		return -1;
	}

	if(CheckBridge(bridge) != 0) {
		return -1;
	}
	
	sock = socket(AF_LOCAL, SOCK_STREAM, 0);
	if(sock <= 0) {
		return -1;
	}

	strncpy(ir.ifr_name, bridge, IFNAMSIZ);
	ir.ifr_ifindex = ifindex;
	if(ioctl(sock, SIOCBRADDIF, &ir) < 0) {
		close(sock);
		return -1;
	}
	close(sock);
	return 0;
}
