#include <iostream>

#include <hurlean/src/client.h>


enum MyMessageType
{
	HELLO
};


int main()
{
	hl::Client<MyMessageType> client;
	
	std::string ip_address;
	unsigned short port;
	
	std::cin >> ip_address;
	std::cin >> port;
	
	if (client.connect(ip_address, port))
	{
		std::cout << "Client has connected!" << std::endl;

		auto& incoming = client.get_incoming_messages();

		while (true)
		{
			if (!incoming.is_empty())
			{
				for (int i = 0; i < incoming.size(); i++)
				{
					auto message = incoming.pop_front();
				}
			}
		}
	}
	else
	{
		std::cout << "Client hasn't been able to connect!" << std::endl;
	}
	
	return 0;
}