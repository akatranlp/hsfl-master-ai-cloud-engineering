# Kubernetes

## Pods

- Eine oder mehrere Container bilden eine app
- wenn der Container stirbt, stirbt der pod
- kleinste einheit

## Deployment

- Überwacht den Lifecycle der zugewiesenen Pods
- wenn ein Pod stirbt wird ein neuer gestartet
- beim update wird darauf geachtet das es nie zu einem kompletten ausfall kommt

## Namespaces

- wie man es kennt kann man hierdurch eine abgrenzung erstellen
- es können network policies gebaut werden, um z.b. namespaces zu isolieren

## Secrets

- Ein sicherer Ort für sensible Daten
- können bei einem Deployment eingebunden werden und dem pod als env mitgegegeben werden

## Services

- Ein Pod ist selbst bei öffnen von Pods nur innerhalb des Clusters zu erreichen über seine spezifische IP-Adresse
- Ein Service ist ein Loadbalancer über eine mehrzahl von Pods mit einem matching label
- meist direkt beim deployment miterstellt
- man erhält eine IP-Adresse die auf allen nodes des Clusters gültigkeit hat, aber immer noch nur intern erreichbar ist

## Ingress

- Hierfür wird ein Ingress-Controller wie nginx oder traefik benötigt
- Er ist die Schnittstelle in die Außenwelt und ein Reverse-Proxy
- Hier werden Pfade und Domains spezifiziert, um dann den oder die zugehörigen Services zu erhalten

## Volumes

- Docker-Volumes sind ja lokal auf einem Rechner verbunden, im Falle von Kubernetes funktioniert es so nicht
- Network attached im ganzen Cluster zur Verfügung
