#include <iostream>

#include <asio.hpp>

#include "server.h"


int main()
{
	enum MessageType
	{
		HELLO
	};
	
	hl::Server<MessageType> server(4545);
	server.start();

	while (true)
	{
		server.update(-1);
	}
	
	return 0;
}