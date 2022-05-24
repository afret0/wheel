#!/user/bin/env bash
# get all filename in specified path

mv_file_to_dir(){
    path=$1
    files=$(ls $path | grep .log)
    for file in $files
    do
        echo $file
#        mv $file $file+".csv"
        rm $file
    done
}
mv_file_to_dir
