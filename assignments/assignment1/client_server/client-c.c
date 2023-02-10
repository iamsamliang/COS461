/*****************************************************************************
 * client-c.c
 * Name: Sam Liang
 * NetId: saml
 *****************************************************************************/

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <errno.h>
#include <string.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <netdb.h>
#include <netinet/in.h>
#include <errno.h>

#define SEND_BUFFER_SIZE 2048

/* TODO: client()
 * Open socket and send message from stdin.
 * Return 0 on success, non-zero on failure
 */
int client(char *server_ip, char *server_port)
{
  int sockfd, numbytes;
  char buffer[SEND_BUFFER_SIZE];
  struct addrinfo hints, *servinfo, *p;
  int rv;
  size_t bytes_read;

  memset(&hints, 0, sizeof hints);
  hints.ai_family = AF_UNSPEC;
  hints.ai_socktype = SOCK_STREAM;

  if ((rv = getaddrinfo(server_ip, server_port, &hints, &servinfo)) != 0)
  {
    fprintf(stderr, "client getaddrinfo error: %s\n", gai_strerror(rv));
    return 1;
  }

  // loop through all the results and connect to the first we can
  for (p = servinfo; p != NULL; p = p->ai_next)
  {
    if ((sockfd = socket(p->ai_family, p->ai_socktype,
                         p->ai_protocol)) == -1)
    {
      perror("client: socket() error");
      continue;
    }

    if (connect(sockfd, p->ai_addr, p->ai_addrlen) == -1)
    {
      close(sockfd);
      perror("client: connect() error");
      continue;
    }

    break;
  }

  if (p == NULL)
  {
    fprintf(stderr, "client: failed to connect to a server\n");
    return 2;
  }

  freeaddrinfo(servinfo); // all done with this structure

  // keep reading bytes from stdin until we have read everything. For each chunk, send it immediately to server
  while ((bytes_read = fread(buffer, sizeof(char), SEND_BUFFER_SIZE, stdin)))
  {
    // fwrite(buffer, 1, bytes_read, stdout);
    // fflush(stdout);
    if (send(sockfd, buffer, bytes_read, 0) == -1)
    {
      perror("error with sending data to server");
      close(sockfd);
      return 2;
    }
  }

  close(sockfd);
  return 0;
}

/*
 * main()
 * Parse command-line arguments and call client function
 */
int main(int argc, char **argv)
{
  char *server_ip;
  char *server_port;

  if (argc != 3)
  {
    fprintf(stderr, "Usage: ./client-c [server IP] [server port] < [message]\n");
    exit(EXIT_FAILURE);
  }

  server_ip = argv[1];
  server_port = argv[2];
  return client(server_ip, server_port);
}
