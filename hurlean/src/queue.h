#pragma once

#include <deque>
#include <mutex>


namespace hl
{
	// thread-safe double-ended queue
	template <class T>
	class Queue
	{
	private:
		std::deque<T> q;
		std::mutex q_mutex;

	public:
		Queue() = default;

		Queue(const Queue<T>&) = delete;

		~Queue()
		{
			clear();
		}

	public:
		// peek at front
		const T& front()
		{
			std::scoped_lock lock(q_mutex);
			return q.front();
		}

		// peek at back
		const T& back()
		{
			std::scoped_lock lock(q_mutex);
			return q.back();
		}

		// remove and return the item at front
		T pop_front()
		{
			std::scoped_lock lock(q_mutex);
			auto item = std::move(q.front());
			q.pop_front();
			return item;
		}

		// remove and return the item at back
		T pop_back()
		{
			std::scoped_lock lock(q_mutex);
			auto item = std::move(q.back());
			q.pop_back();
			return item;
		}

		// push an item to the back
		void push_back(const T& item)
		{
			std::scoped_lock lock(q_mutex);
			q.emplace_back(std::move(item));
		}

		// push an item to the back
		void push_front(const T& item)
		{
			std::scoped_lock lock(q_mutex);
			q.emplace_front(std::move(item));
		}

		// returns true if queue is empty
		bool is_empty()
		{
			std::scoped_lock lock(q_mutex);
			return q.empty();
		}

		// returns number of items in queue
		size_t size()
		{
			std::scoped_lock lock(q_mutex);
			return q.size();
		}

		// clear queue
		void clear()
		{
			std::scoped_lock lock(q_mutex);
			q.clear();
		}
	};
}