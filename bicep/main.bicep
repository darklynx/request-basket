/*

az login
az deployment group create --resource-group ... --template-file main.bicep --parameters basketAdminPassword=yolo

*/

param location string = 'westeurope'
param projectName string = 'helloAzure'
param basketAdminToken string

resource appServicePlan 'Microsoft.Web/serverfarms@2023-12-01' = {
  name: 'plan-baskets-${projectName}-001'
  location: location
  kind: 'linux'
  properties: {
    reserved: true
  }
  sku: {
    name: 'F1'
    tier: 'Free'
  }
}

resource webApp 'Microsoft.Web/sites@2023-12-01' = {
  name: 'app-baskets-${projectName}-001'
  location: location
  properties: {
    serverFarmId: appServicePlan.id
    siteConfig: {
      linuxFxVersion: 'DOCKER|index.docker.io/darklynx/request-baskets'
      appSettings: [
        {
          name: 'DOCKER_REGISTRY_SERVER_PASSWORD'
          value: 'xxxxx'
        }
        {
          name: 'DOCKER_REGISTRY_SERVER_URL'
          value: 'https://index.docker.io'
        }
        //following settings are for the container itself, please see: https://github.com/darklynx/request-baskets?tab=readme-ov-file#parameters
        {
          name: 'MAXSIZE' //maximum allowed basket capacity, basket capacity greater than this number will be rejected by service
          value: '77777'
        }
        {
          name: 'TOKEN' // master token to gain control over all baskets, if not defined a random token will be generated when service is launched and printed to stdout
          value: basketAdminToken
        }
        {
          name: 'BASKET' //name of a basket to auto-create during service startup, this parameter can be specified multiple times
          value: 'helloBasket'
        }
        {
          name: 'MODE' //defines service operation mode: public - when any visitor can create a new basket, or restricted - baskets creation requires master token
          value: 'public'
        }
        {
          name: 'THEME' //CSS theme for web UI, supported values: standard, adaptive, flatly
          value: 'flatly'
        }
        {
          name: 'PAGE' //default page size when retrieving collections
          value: '77777'
        }
      ]
    }
  }
}
