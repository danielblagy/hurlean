#pragma once

#include <thread>
#include <memory>

#include "connection.h"


namespace hl
{
	// Manages the connection object (manages in/out messages, provides a thread to listen on connection)
	template <class T>
	class ClientSession
	{
	private:
		std::shared_ptr<Connection<T>> connection;
		std::thread read_thread;
		std::thread write_thread;
		bool running;

	public:
		ClientSession(std::shared_ptr<Connection<T>> _connection)
			: connection(std::move(_connection))
		{
			// start the session
			running = true;
			read_thread = std::thread(&ClientSession::listen, this);
			write_thread = std::thread(&ClientSession::send, this);
		}

		// a move constructor, because read_thread and write_thread are not copyable
		ClientSession(ClientSession&& o)
			: read_thread(std::move(o.read_thread)), write_thread(std::move(o.write_thread))
		{}

		~ClientSession()
		{
			running = false;
			read_thread.join();
			write_thread.join();
		}

	public:
		void send_to_client(const Message<T>& message)
		{
			connection->out.push_back(std::move(message));
		}

		Queue<Message<T>> get_incoming_messages()
		{
			return connection->in;
		}

	private:
		// used as an update function in connection thread for listening for messages from the client
		void listen()
		{
			while ((running = connection->opened))
			{
				if (connection->read())
				{
					connection->in.push_back(connection->temp_message_in);
				}
				else
				{
					connection->close();
				}
			}
		}

		// used as an update function in connection thread for listening for messages from the client
		void send()
		{
			while ((running = connection->opened))
			{
				if (connection->write(connection->out.pop_front()))
				{
					// TODO : maybe log
				}
				else
				{
					connection->close();
				}
			}
		}
	};
}