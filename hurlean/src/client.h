#pragma once

#include <string>
#include <memory>
#include <thread>

#include <asio.hpp>

#include "connection.h"
#include "session.h"


namespace hl
{
	template <class T>
	class Client
	{
	private:
		asio::io_context io_context;

		std::unique_ptr<Session<T>> server_session;

	public:
		// port on which server is going to be running
		Client()
		{}

		~Client()
		{
			disconnect();
			server_session.release();
		}

		bool connect(const std::string& host, unsigned short port)
		{
			try
			{
				// Resolve hostname/ip-address into tangiable physical address
				asio::ip::tcp::resolver resolver(io_context);
				asio::ip::tcp::resolver::results_type endpoints = resolver.resolve(host, std::to_string(port));

				std::shared_ptr<Connection<T>> server_connection = std::make_shared<Connection<T>>(io_context);

				if (!server_connection->connect_to_server(endpoints))
				{
					// TODO : handle the server_session object
					return false;
				}

				server_session = std::make_unique<Session<T>>(server_connection);
			}
			catch (std::exception& e)
			{
				std::cerr << "Client Exception: " << e.what() << "\n";

				return false;
			}

			return true;
		}

		void disconnect()
		{
			server_session->close();
		}

		bool connected()
		{
			return server_session->is_open();
		}

		void send(const Message<T>& message)
		{
			if (connected())
				server_session->send(message);
		}

		Queue<Message<T>>& get_incoming_messages()
		{
			return server_session->get_incoming_messages();
		}
	};
}