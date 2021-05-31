# chart-prettier

Chart prettier reads the template folder of a helm chart and places each object in its own file.

In case multiple manifests of the same type exist, the name schema {type}-{name}.yaml is used for this type.

## Example

``` bash
chart-prettier chart/templates
```

``` bash
cat all-the-things.yaml | chart-prettier chart/templates
```
