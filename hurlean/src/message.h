#pragma once

#include <vector>


namespace hl
{
	template <class T>
	class Message {
		
		T type;
		uint32_t size;

		std::vector<uint8_t> data;

	};
}