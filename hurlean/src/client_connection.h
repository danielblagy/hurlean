#pragma once

#include <thread>

#include <asio.hpp>


namespace hl
{
	class ClientConnection
	{
	private:
		asio::ip::tcp::socket socket;
		std::thread connection_thread;
	};
}