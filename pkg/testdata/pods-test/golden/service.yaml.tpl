apiVersion: v1
kind: Service
metadata:
  name: my-nginx
  labels:
    run: my-nginx
spec:
  ports:
  - port: {{ .Values.Port }}
    protocol: TCP
  selector:
    run: my-nginx