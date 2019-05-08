#!/bin/bash
# get all filename in specified path

mv_file_to_dir(){
    path=$1
    files=$(ls $path | grep -v _1.csv)
    for file in $files
    do
        echo $file
        dir=`echo $file | cut -d \. -f 1`
        mkdir $dir
        mv $file $dir
    done
}
mv_file_to_dir
