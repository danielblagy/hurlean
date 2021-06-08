#pragma once

#include <string>
#include <vector>
#include <memory>
#include <thread>

#include <asio.hpp>

#include "connection.h"
#include "client_session.h"


namespace hl
{
	template <class T>
	class Server
	{
	private:
		asio::io_context io_context;
		asio::ip::tcp::acceptor acceptor;
		
		std::vector<ClientSession<T>> client_sessions;

		bool running;
		std::thread accept_thread;

	public:
		// port on which server is going to be running
		Server(unsigned short port)
			: acceptor(io_context, asio::ip::tcp::endpoint(asio::ip::tcp::v4(), port))
		{}
		
		~Server()
		{
			running = false;
			accept_thread.join();
		}

		void start()
		{
			running = true;
			accept_thread = std::thread(&Server::accept_client_connections, this);
		}

		void update()
		{
			
		}
	
	public:
		virtual void on_client_connect() {};
		virtual void on_client_disconnect() {};
		virtual void on_client_message() {};

	private:
		void accept_client_connections()
		{
			while (running)
			{
				std::shared_ptr<Connection<T>> client_connection = std::make_shared<Connection<T>>(io_context);
				client_connection->wait(acceptor);

				// we got connection

				client_sessions.emplace_back(client_connection);
				ClientSession& client_session = client_sessions.back();

				std::string message = "hello from the class";
				client_session.send_to_client(Message<T>( {} ));
			}
		}
		
		void write()
		{}
	};
}