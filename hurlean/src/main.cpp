#include <iostream>

#include <asio.hpp>

#include "server.h"


enum MyMessageType
{
	HELLO
};

class MyServer : public hl::Server<MyMessageType>
{
public:
	MyServer(unsigned short port)
		: Server(port)
	{}

protected:
	virtual void on_client_connect(hl::ClientSession<MyMessageType>& client) override
	{
		std::string message_string = "hello from the class";
		int data = 15;

		hl::Message<MyMessageType> hello_message;
		hello_message.header.type = MyMessageType::HELLO;
		//hello_message.header.size = sizeof(data);
		hello_message << data;

		client.send_to_client(hello_message);

	}

	virtual void on_client_disconnect(hl::ClientSession<MyMessageType>& client) override
	{}

	virtual void on_client_message(hl::ClientSession<MyMessageType>& client, const hl::Message<MyMessageType>& message) override
	{}
};


int main()
{
	MyServer server(4545);
	server.start();

	while (true)
	{
		server.update();
	}
	
	return 0;
}