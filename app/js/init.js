var welcome ="e=>end:>http://www.google.com\n"+
"op1=>operation: My Operation\n"+
"sub1=>subroutine: My Subroutine\n"+
"cond=>condition: Yes\n"+
"or No?:>http://www.google.com\n"+
"io=>inputoutput: catch something...\n"+
"para=>parallel: parallel tasks\n"+
"st->op1->cond\n"+
"cond(yes)->io->e\n"+
"cond(no)->para\n"+
"para(path1, bottom)->sub1(right)->op1\n"+
"para(path2, top)->op1\n"+
"st=>start: Start:>http://www.google.com[blank]\n";