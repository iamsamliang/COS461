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
  // struct sigaction sa;
  struct sockaddr_storage their_addr;
  socklen_t sin_size;
  int yes = 1;
  char buffer[RECV_BUFFER_SIZE];
  int num_rec_bytes;
  char s[INET6_ADDRSTRLEN];

  memset(&hints, 0, sizeof hints); // make sure the struct is empty
  hints.ai_family = AF_UNSPEC;     // don't care IPv4 or IPv6
  hints.ai_socktype = SOCK_STREAM; // TCP stream sockets
  hints.ai_flags = AI_PASSIVE;     // use my own IP for me (localhost)

  if ((status = getaddrinfo(NULL, server_port, &hints, &servinfo)) != 0)
  {
    fprintf(stderr, "server getaddrinfo error: %s\n", gai_strerror(status));
    return 1;
  }

  // servinfo now points to a linked list of 1 or more struct addrinfos

  // loop through all the results and bind to the first we can
  for (p = servinfo; p != NULL; p = p->ai_next)
  {
    if ((sockfd = socket(p->ai_family, p->ai_socktype,
                         p->ai_protocol)) == -1)
    {
      perror("server: error with using this socket");
      continue;
    }

    if (setsockopt(sockfd, SOL_SOCKET, SO_REUSEADDR, &yes,
                   sizeof(int)) == -1)
    {
      perror("server: error setting this socket option");
      return 1;
    }

    if (bind(sockfd, p->ai_addr, p->ai_addrlen) == -1)
    {
      close(sockfd);
      perror("server: error binding to this socket");
      continue;
    }

    break;
  }

  freeaddrinfo(servinfo); // free the linked-list

  if (p == NULL)
  {
    fprintf(stderr, "server: failed to bind to any socket\n");
    return 1;
  }

  // listen for connections
  if (listen(sockfd, QUEUE_LENGTH) == -1)
  {
    perror("error with trying to listen for client connections");
    return 1;
  }

  while (1)
  {
    // accept connections while continuing to listen for more connections and adding them to queue
    sin_size = sizeof their_addr;
    new_fd = accept(sockfd, (struct sockaddr *)&their_addr, &sin_size);

    if (new_fd == -1)
    {
      perror("error with accepting client connection");
      continue;
    }

    // write to stdout as we receive the messages
    while (1)
    {
      num_rec_bytes = recv(new_fd, buffer, RECV_BUFFER_SIZE, 0);

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
        // buffer[num_rec_bytes] = '\0';
        // printf("%s", buffer);
        // write(new_fd, buffer, num_rec_bytes);
        fwrite(buffer, 1, num_rec_bytes, stdout);
        fflush(stdout);
        // fprintf(stdout, "%s", buffer);
      }
    }

    // close connection to current client when done
    close(new_fd);
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
