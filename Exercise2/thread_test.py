import threading

i = 0
lock = threading.Lock()

def increase():
	global i

	for j in range(100000):
		lock.acquire()
		i += 1
		lock.release()

def decrease():
	global i

	for j in range(100000):
		lock.acquire()
		i -= 1
		lock.release()

def main():
	thread1 = threading.Thread(target = increase, args = (),)
	thread2 = threading.Thread(target = decrease, args = (),)

	thread1.start()
	thread2.start()

	thread1.join()
	thread2.join()

	print i

main()
