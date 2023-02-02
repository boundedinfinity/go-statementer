#!/usr/bin/env fish

function ee; echo $argv; eval $argv; end

set script_dir (cd (dirname (status -f)); and pwd)
set pdf_path $argv[1]
set pdf_fullname (basename $pdf_path)
set pdf_nameonly (echo $pdf_fullname | cut -f 1 -d '.')
set pdf_workdir $script_dir/$pdf_nameonly
set IMAGE_EXT png

# set fish_trace 1

# echo "          Source: $pdf_path"
# echo "  Work Directory: $pdf_workdir"

if test -d $pdf_workdir
    rm -rf $pdf_workdir/
end

mkdir -p $pdf_workdir
cp $pdf_path $pdf_workdir/$pdf_fullname

./convert.fish $pdf_workdir -density 300 -quiet $pdf_fullname $pdf_nameonly-%04d.$IMAGE_EXT

for image_path in $pdf_workdir/*.$IMAGE_EXT
    set image_fullname (basename $image_path)
    set image_onlyname (echo $image_fullname | cut -f 1 -d '.')
    ./tesseract.fish $pdf_workdir --oem 1 -l eng --psm 6 -c preserve_interword_spaces=1 $image_fullname $image_onlyname
end

for text_path in $pdf_workdir/*.txt
    set text_fullname (basename $text_path)
    cat $text_path >> $pdf_workdir/$pdf_nameonly.txt
end