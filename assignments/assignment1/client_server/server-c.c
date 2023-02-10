/*****************************************************************************
 * server-c.c
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

#define QUEUE_LENGTH 10
#define RECV_BUFFER_SIZE 2048

void sigchld_handler(int s)
{
  // waitpid() might overwrite errno, so we save and restore it:
  int saved_errno = errno;

  while (waitpid(-1, NULL, WNOHANG) > 0)
    ;

  errno = saved_errno;
}

/* TODO: server()
 * Open socket and wait for client to connect
 * Print received message to stdout
 * Return 0 on success, non-zero on failure
 */
int server(char *server_port)
{
  int status;
  struct addrinfo hints;     // specifications for server from getaddrinfo call
  struct addrinfo *servinfo; // will point to the results
  struct addrinfo *p;
  int sockfd; // socket to listen on
  int new_fd; // accept returns new socket descriptor each time while old one keeps listening
  struct sigaction sa;
  struct sockaddr_storage their_addr;
  socklen_t sin_size;
  int yes = 1;
  char buffer[RECV_BUFFER_SIZE];
  int num_rec_bytes;

  memset(&hints, 0, sizeof hints); // make sure the struct is empty
  hints.ai_family = AF_UNSPEC;     // don't care IPv4 or IPv6
  hints.ai_socktype = SOCK_STREAM; // TCP stream sockets
  hints.ai_flags = AI_PASSIVE;     // use my own IP for me (localhost)

  if ((status = getaddrinfo(NULL, server_port, &hints, &servinfo)) != 0)
  {
    fprintf(stderr, "getaddrinfo error: %s\n", gai_strerror(status));
    return 1;
  }

  // servinfo now points to a linked list of 1 or more struct addrinfos

  // loop through all the results and bind to the first we can
  for (p = servinfo; p != NULL; p = p->ai_next)
  {
    if ((sockfd = socket(p->ai_family, p->ai_socktype,
                         p->ai_protocol)) == -1)
    {
      perror("server: socket");
      continue;
    }

    if (setsockopt(sockfd, SOL_SOCKET, SO_REUSEADDR, &yes,
                   sizeof(int)) == -1)
    {
      perror("setsockopt");
      return 1;
    }

    if (bind(sockfd, p->ai_addr, p->ai_addrlen) == -1)
    {
      close(sockfd);
      perror("server: bind");
      continue;
    }

    break;
  }

  freeaddrinfo(servinfo); // free the linked-list

  if (p == NULL)
  {
    fprintf(stderr, "server: failed to bind\n");
    return 1;
  }

  if (listen(sockfd, QUEUE_LENGTH) == -1)
  {
    perror("listen");
    return 1;
  }

  sa.sa_handler = sigchld_handler; // reap all dead processes
  sigemptyset(&sa.sa_mask);
  sa.sa_flags = SA_RESTART;
  if (sigaction(SIGCHLD, &sa, NULL) == -1)
  {
    perror("sigaction");
    return 1;
  }

  printf("server: waiting for connections...\n");

  while (1)
  { // main accept() loop
    sin_size = sizeof their_addr;
    new_fd = accept(sockfd, (struct sockaddr *)&their_addr, &sin_size);

    if (new_fd == -1)
    {
      perror("error with accepting client connection");
      continue;
    }

    // this prints client information
    // inet_ntop(their_addr.ss_family,
    //           get_in_addr((struct sockaddr *)&their_addr),
    //           s, sizeof s);
    // printf("server: got connection from %s\n", s);

    while (1)
    {
      memset(buffer, 0, RECV_BUFFER_SIZE);
      num_rec_bytes = recv(new_fd, buffer, RECV_BUFFER_SIZE, 0);

      // TODO: should this return failure???
      if (num_rec_bytes == -1)
      {
        perror("error with receiving client data");
      }
      else if (num_rec_bytes == 0)
      {
        break;
      }
      else
      {
        fprintf(stdout, "%s", buffer);
      }
    }

    // if ((numbytes = recv(sockfd, buffer, SEND_BUFFER_SIZE - 1, 0)) == -1)
    // {
    //   perror("recv");
    //   exit(1);
    // }

    // buf[numbytes] = '\0';

    // do not need to fork new process for this assignment
    // if (!fork())
    // {                // this is the child process
    //   close(sockfd); // child doesn't need the listener
    //   if (recv(new_fd, "Hello, world!", RECV_BUFFER_SIZE, 0) == -1)
    //     perror("send");
    //   close(new_fd);
    //   exit(0);
    // }

    close(new_fd); // parent doesn't need this
  }

  return 0;
}

/*
 * main():
 * Parse command-line arguments and call server function
 */
int main(int argc, char **argv)
{
  char *server_port;

  if (argc != 2)
  {
    fprintf(stderr, "Usage: ./server-c [server port]\n");
    exit(EXIT_FAILURE);
  }

  server_port = argv[1];
  return server(server_port);
}
