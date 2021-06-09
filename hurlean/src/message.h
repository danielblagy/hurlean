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

		// Pushes any POD-like data into the message buffer
		template<typename DataType>
		friend Message<T>& operator << (Message<T>& msg, const DataType& data)
		{
			// Check that the type of the data being pushed is trivially copyable
			static_assert(std::is_standard_layout<DataType>::value, "Data is too complex to be pushed into vector");

			// Cache current size of vector, as this will be the point we insert the data
			size_t previous_size = msg.body.size();

			// Resize the vector by the size of the data being pushed
			msg.body.resize(msg.body.size() + sizeof(DataType));

			// Physically copy the data into the newly allocated vector space
			std::memcpy(msg.body.data() + previous_size, &data, sizeof(DataType));

			// Recalculate the message size
			msg.header.size = msg.body.size();

			// Return the target message so it can be "chained"
			return msg;
		}

		// Pulls any POD-like data form the message buffer
		template<typename DataType>
		friend Message<T>& operator >> (Message<T>& msg, DataType& data)
		{
			// Check that the type of the data being pushed is trivially copyable
			static_assert(std::is_standard_layout<DataType>::value, "Data is too complex to be pulled from vector");

			// Cache the location towards the end of the vector where the pulled data starts
			size_t start_location = msg.body.size() - sizeof(DataType);

			// Physically copy the data from the vector into the user variable
			std::memcpy(&data, msg.body.data() + start_location, sizeof(DataType));

			// Shrink the vector to remove read bytes, and reset end position
			msg.body.resize(start_location);

			// Recalculate the message size
			msg.header.size = msg.body.size();

			// Return the target message so it can be "chained"
			return msg;
		}
	};
}