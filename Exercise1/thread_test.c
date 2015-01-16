#include <pthread.h>
#include <stdio.h>

void * increase(void);
void * decrease(void);

int main(void)
{
	static int i = 0;

	pthread_t thread_1;
	pthread_t thread_2;

	pthread_create(&thread_1, NULL, increase, NULL);
	pthread_create(&thread_2, NULL, decrease, NULL);

	pthread_join(thread_1, NULL);
	pthread_join(thread_2, NULL);

	printf("%d\n", i);

	return 0;
}

void * increase(void)
{
	static int i;

	for (int j = 0; j < 1000000; ++j)
	{
		i++;
	}
}

void * decrease(void)
{
	static int i;

	for (int j = 0; j < 1000000; ++j)
	{
		i--;
	}
}