from threading import Thread

i = 0

def main():
	thread1 = Thread(target = thread1foo, args = (),)
	thread1.start()

	thread2 = Thread(target = thread2foo, args = (),)
	thread2.start()

	thread1.join()
	thread2.join()

	print i

def thread1foo():
	global i

	for j in range(1000000):
		i += 1

def thread2foo():
	global i

	for j in range(1000000):
		i -= 1

main()