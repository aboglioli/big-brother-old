# TODO

x Improve Error handling: validation and internal errors
x Update uses asynchronously:
    x Add 'usesUpdated' property to Composition, default: false
    x Emit event of 'CompositionUpdatedManually' after updating a Composition
    x Make a service to react to that event, update uses and publish 'CompositionUpdatedAutomatically'
    x Set 'usesUpdated' from Composition updated manually to true