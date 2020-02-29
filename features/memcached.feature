Feature: Memcached deployment
  As a Kubernetes
  I want a cloud native memcached application
  So that I can deploy a managed instance of memcached in my namespace
  
  Acceptance Criteria
  - Deploys the desired number of memcached instances
  - Creates a memcached service
  - Cleans up all created resources when CR is deleted

  Scenario: Creating a memcached
    Given I create a Resource:
    """yaml
    apiVersion: cache.example.com/v1alpha1
    kind: Memcached
    metadata:
      name: example-memcached
      namespace: default
    spec:
      size: 3
    """
    When the operator reconciles
    Then there should exist a "v1" "" "Service" called "example-memcached" in namespace "default"
    And there should exist a "v1" "apps" "Deployment" called "example-memcached" in namespace "default"
    And there should exist 3 ready pods for Deployment called "example-memcached" in namespace "default"
    And there should exist 3 node names in the status of Memcached "example-memcached" in namespace "default"

  Scenario: Creating a memcached
    Given I create a Resource:
    """yaml
    apiVersion: cache.example.com/v1alpha1
    kind: Memcached
    metadata:
      name: another-memcached
      namespace: default
    spec:
      size: 5
    """
    When the operator reconciles
    Then there should exist a "v1" "" "Service" called "another-memcached" in namespace "default"
    And there should exist a "v1" "apps" "Deployment" called "another-memcached" in namespace "default"
    And there should exist 5 ready pods for Deployment called "another-memcached" in namespace "default"
    And there should exist 5 node names in the status of Memcached "another-memcached" in namespace "default"

  Scenario: Creating and Deleting memcached
    Given I create a Resource:
    """yaml
    apiVersion: cache.example.com/v1alpha1
    kind: Memcached
    metadata:
      name: example-memcached
      namespace: default
    spec:
      size: 3
    """
    When the operator reconciles
    Then there should exist a "v1" "" "Service" called "example-memcached" in namespace "default"
    And there should exist a "v1" "apps" "Deployment" called "example-memcached" in namespace "default"
    When I delete a "v1alpha1" "cache.example.com" "Memcached" called "example-memcached" in namespace "default"
    And the operator reconciles
    Then there should not exist a "v1" "" "Service" called "example-memcached" in namespace "default"
    And there should not exist a "v1" "apps" "Deployment" called "example-memcached" in namespace "default"

  Scenario: Reconciling a controlled Deployment
    Given I create a Resource:
    """yaml
    apiVersion: cache.example.com/v1alpha1
    kind: Memcached
    metadata:
      name: example-memcached
      namespace: default
    spec:
      size: 1
    """
    When the operator reconciles
    Then there should exist a "v1" "apps" "Deployment" called "example-memcached" in namespace "default"
    When I delete a "v1" "apps" "Deployment" called "example-memcached" in namespace "default"
    And the operator reconciles
    Then there should exist a "v1" "apps" "Deployment" called "example-memcached" in namespace "default"


  Scenario: Reconciling a controlled Service
    Given I create a Resource:
    """yaml
    apiVersion: cache.example.com/v1alpha1
    kind: Memcached
    metadata:
      name: example-memcached
      namespace: default
    spec:
      size: 3
    """
    When the operator reconciles
    Then there should exist a "v1" "" "Service" called "example-memcached" in namespace "default"
    When I delete a "v1" "" "Service" called "example-memcached" in namespace "default"
    And the operator reconciles
    Then there should exist a "v1" "" "Service" called "example-memcached" in namespace "default"

