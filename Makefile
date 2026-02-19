all:
	c-for-go -ccdefs -out pkg/ spirv-reflect.yml

clean:
	rm -f pkg/spvreflect/*

test:
	cd pkg/spvreflect/ && go build