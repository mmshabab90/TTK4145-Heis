#include <pthread.h>
#include <stdio.h>

void thread_1_foo(void);
void thread_2_foo(void);

int main(void)
{
	static int i = 0;

	pthread_t thread_1;
	pthread_create(&thread_1, NULL, thread_1_foo, NULL);

	pthread_t thread_2;
	pthread_create(&thread_2, NULL, thread_2_foo, NULL);

	pthread_join(thread_1, NULL);
	pthread_join(thread_2, NULL);

	printf("%d\n Oyvind sugar", i);

	return 0;
}

void thread_1_foo(void)
{
	static int i;

	for (int j = 0; j < 1000000; ++j)
	{
		i++;
	}
}

void thread_2_foo(void)
{
	static int i;
	
	for (int j = 0; j < 1000000; ++j)
	{
		i--;
	}
}