#pragma once

#include <deque>

#include <asio.hpp>

#include "connection.h"


namespace hl
{
	template <class T>
	class Connection
	{
	private:
		asio::ip::tcp::socket socket;
		std::deque<Message> in;
		std::deque<Message> out;

	public:
		Connection(const asio::io_context& io_context);
		~Connection() = default;

	};
}