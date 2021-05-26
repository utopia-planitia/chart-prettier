# csort

`csort` stands for chart sorter.

Chart sorter reads the template folder of a helm chart and places each object in its own file.

In case multiple manifests of the same type exist, the name schema for this type switches to {type}-{name}.yaml.

## Example

``` bash
csort chart/templates
```

``` bash
cat all-the-things.yaml | csort chart/templates
```
