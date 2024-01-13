# Kubernetes

This file is in german and describes the concepts of Kubernetes

## Pods

- Eine oder mehrere Container bilden eine app
- Wenn der Container stirbt, stirbt der pod
- Kleinste einheit

## Deployment

- Überwacht den Lifecycle der zugewiesenen Pods
- Wenn ein Pod stirbt wird ein neuer gestartet
- Beim Update wird darauf geachtet das es nie zu einem kompletten Ausfall kommt

## Namespaces

- Wie man es kennt, kann man hierdurch eine Abgrenzung erstellen
- Es können network policies gebaut werden, um z.b. namespaces zu isolieren

## Secrets

- Ein sicherer Ort für sensible Daten
- können bei einem Deployment eingebunden werden und dem Pod als env mitgegegeben werden

## Services

- Ein Pod ist selbst beim Öffnen von Pods nur innerhalb des Clusters zu erreichen, über seine spezifische IP-Adresse
- Ein Service ist ein Loadbalancer über eine mehrzahl von Pods mit einem matching label
- Meist direkt beim Deployment miterstellt
- Man erhält eine IP-Adresse die auf allen nodes des Clusters gültigkeit hat, aber immer noch nur intern erreichbar ist

## Ingress

- Hierfür wird ein Ingress-Controller wie nginx oder traefik benötigt
- Er ist die Schnittstelle in die Außenwelt und ein Reverse-Proxy
- Hier werden Pfade und Domains spezifiziert, um dann den oder die zugehörigen Services zu erhalten

## Volumes

- Docker-Volumes sind ja lokal auf einem Rechner verbunden, im Falle von Kubernetes funktioniert es so nicht
- Network attached im ganzen Cluster zur Verfügung
