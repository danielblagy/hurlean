#pragma once

#include <thread>
#include <string>
#include <vector>

#include <asio.hpp>


namespace hl
{
	class Server
	{
	private:
		std::vector<std::thread> threads;
	};
}