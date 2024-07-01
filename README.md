# [WIP] Polybar module for bspwm windows list with support multiple monitors

```bash
# build
make build

# install
mkdir -p ~/.config/polybar/scripts/bspwm-windows/bspcw
cp bspcw ~/.config/polybar/scripts/bspwm-windows/bspcw
cp config.toml ~/.config/polybar/scripts/bspwm-windows/config.toml
```

### Module
```ini 
[module/bspwm-windows]
type = custom/script
exec = ~/.config/polybar/scripts/bspwm-windows/bspcw watch -m $MONITOR --config ~/.config/polybar/scripts/bspwm-windows/config.toml
tail = true
format = <label>
# format-font = 2
# separator = " "
# label-padding = 1
```

If you have single monitor you can remove `-m $MONITOR` from exec

// TODO

Based on [diogox/bspc-go](https://github.com/diogox/bspc-go) and [tuurep/windowlist](https://github.com/tuurep/windowlist)


