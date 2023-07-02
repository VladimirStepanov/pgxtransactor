What is that?
-------------

**pgxtransactgor** is a library for using a common transaction connection between repository method calls


How to use it?
--------------

```go
//...
txManager.RunRepeatableRead(ctx, func(dbCtx context.Context) error {
    product, err := repo.GetProduct(dbCtx, 1)
    if err != nil {
        return err
    }
    // some data processing
    product.Price += 10000
    return repo.UpdateProduct(dbCtx, product)
})
```

Full example in `./example` directory
