# Web-Service
For the web service, we utilize Elm and vite <3.

## Development 

Install the node depnendcies with yarn
```
yarn
```

Generate the tailwind classes for elm with:

```
yarn generate
```

Start the vite development server with:
```
yarn dev
```

The requests to other services will be proxied by the `vite.config.cjs`. 
You can just use the default docker-compose configuration for that.

## Prodcution

To build and minify the assets run `yarn build`
```
yarn build
```
