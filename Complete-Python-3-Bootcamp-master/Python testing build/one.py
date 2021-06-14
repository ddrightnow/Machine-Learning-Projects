#one.py
'''
print("hello")
def...



if __name__ ==  "__main__":
	myfunc()
	'''
def func():
	print("FUNC() IN ONE.PY")

def function():
	pass
def function2():
	pass

print("TOP LEVEL IN ONE.PY")

if __name__ == '__main__':
	function()
	function2()

	#print('ONE.PY is being run directly')
#else:
	#print('ONE.PY has been imported')