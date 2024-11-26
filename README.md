# lsproxy

Simple proxy to convert any language server from stdio to TCP and print the
messages being sent between the client and server.

## Use

Adjust the source code as you want with whatever port, working directory,
LSP command, etc. and start the server

```
go run .
```

And then configure your editor to use the proxy instead of starting the LSP.
In Neovim I do it like this:

```
lspconfig.gopls.setup({
    cmd = vim.lsp.rpc.connect("127.0.0.1", 1337)
})
```
