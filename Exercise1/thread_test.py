from threading import Thread

i = 0

def increase():
	global i

	for j in range(1000000):
		i += 1

def decrease():
	global i

	for j in range(1000000):
		i -= 1

def main():
	thread1 = Thread(target = increase, args = (),)
	thread2 = Thread(target = decrease, args = (),)
	
	thread1.start()
	thread2.start()

	thread1.join()
	thread2.join()

	print i

main()