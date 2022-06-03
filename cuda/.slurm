#!/bin/bash

#SBATCH --job-name=
#SBATCH --partition=
#SBATCH --output=
#SBATCH --error=
#SBATCH -N 
#SBATCH --ntasks-per-node=
#SBATCH --cpus-per-task=
#SBATCH --gres=gpu:

ulimit -s unlimited
ulimit -l unlimited

module load gcc/8.3.0 cuda/10.1.243-gcc-8.3.0
nvcc .cu -o  -lcublas

./