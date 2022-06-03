#include <stdio.h>
#include <stdlib.h>
#include <cuda_runtime.h>
 
 
__global__ void MatMul(int *M,int *N,int *res,int width)
{
	int x = threadIdx.x;
	int y = threadIdx.y;

	float elem1 = 0.0,elem2 = 0.0,value = 0.0;
	for(int i = 0;i < width;i++)
	{
		elem1 = M[y * width + i];//取M矩阵的一行
		elem2 = N[i * width + x];//取N矩阵的一列

		value += elem1 * elem2;//求和
	}

	res[y * width + x] = value;
}

__global__ void MatAdd(int *M,int *N,int *res,int width)
{
	int x = threadIdx.x;
	int y = threadIdx.y;

    int value = M[x * width + y] + N[x * width + y];

	res[x * width + y] = value;
}
 
int main()
{
	const int ND = 30;
	int a[ND][ND],b[ND][ND],c[ND][ND],d[ND][ND];
	int *M,*N,*P,*Q;
	
	int width = ND;
	dim3 blockSize(ND,ND);
	
	//设备端内存分配
	cudaMalloc((void**)&M,ND * ND * sizeof(int));
	cudaMalloc((void**)&N,ND * ND * sizeof(int));
	cudaMalloc((void**)&P,ND * ND * sizeof(int));
	cudaMalloc((void**)&Q,ND * ND * sizeof(int));
	
	//初始化
	for(int i = 0;i < ND;i++)
	{
		for(int j = 0;j < ND;j++)
		{
			a[i][j] = 2;
			b[i][j] = 3;
		}
	}
	
	int Size = ND * ND;
	//数据拷贝，主机到设备
	cudaMemcpy(M,a,Size * sizeof(int),cudaMemcpyHostToDevice);
	cudaMemcpy(N,b,Size * sizeof(int),cudaMemcpyHostToDevice);

	MatMul<<<1,blockSize>>>(M,N,P,width);//调用核函数
	MatAdd<<<1,blockSize>>>(M,N,Q,width);//调用核函数
	
	cudaMemcpy(c,P,Size * sizeof(int),cudaMemcpyDeviceToHost);
	cudaMemcpy(d,Q,Size * sizeof(int),cudaMemcpyDeviceToHost);
	
	printf("c0 = %d \n",c[0][0]);
	printf("d0 = %d \n",d[0][0]);
	
	//释放设备内存
	cudaFree(M);
	cudaFree(N);
	cudaFree(P);
	cudaFree(Q);
	
	return 0;
}