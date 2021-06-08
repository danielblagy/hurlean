#pragma once

#include <vector>


namespace hl
{
	template <class T>
	struct MessageHeader
	{
		T type;
		uint32_t size;
	};
	
	template <class T>
	struct Message
	{	
		MessageHeader<T> header;
		std::vector<uint8_t> body;
	};
}