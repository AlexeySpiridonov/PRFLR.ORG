#!/bin/bash

ansible-playbook -i ansible/hosts ansible/playbooks/prod/node_distribute.yml
