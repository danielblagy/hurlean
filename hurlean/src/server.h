#pragma once

#include <string>
#include <vector>
#include <memory>
#include <thread>

#include <asio.hpp>

#include "connection.h"
#include "session.h"


namespace hl
{
	template <class T>
	class Server
	{
	private:
		asio::io_context io_context;
		asio::ip::tcp::acceptor acceptor;
		
		std::vector<Session<T>> client_sessions;

		bool running;
		std::thread accept_thread;

	public:
		// port on which server is going to be running
		Server(unsigned short port)
			: acceptor(io_context, asio::ip::tcp::endpoint(asio::ip::tcp::v4(), port))
		{}
		
		~Server()
		{
			stop();
		}

		void start()
		{
			running = true;
			accept_thread = std::thread(&Server::accept_client_connections, this);
		}

		void stop()
		{
			//io_context.stop();	// I think it's just for async
			running = false;
			accept_thread.join();
		}

		void update()
		{
			
		}
	
	protected:
		virtual void on_client_connect(Session<T>& client) {};
		// TODO : call on diconnect
		virtual void on_client_disconnect(Session<T>& client) {};
		// TODO : call on message
		virtual void on_client_message(Session<T>& client, const Message<T>& message) {};

	private:
		void accept_client_connections()
		{
			while (running)
			{
				std::shared_ptr<Connection<T>> client_connection = std::make_shared<Connection<T>>(io_context);
				client_connection->wait_for_client(acceptor);

				// we got connection

				client_sessions.emplace_back(client_connection);
				//Session& client_session = client_sessions.back();

				on_client_connect(client_sessions.back());
			}
		}
		
		void write()
		{}
	};
}