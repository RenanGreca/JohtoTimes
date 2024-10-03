## Build/run:

1. Ensure GOPATH is in your PATH. Add the following to your .bashrc/.zshrc:
```sh
export PATH="$(go env GOPATH)/bin:$PATH"
```

2. Install dependencies:
```sh
go install github.com/a-h/templ
go install github.com/air-verse/air
```

3. Use air to build:
```sh
air [-c .air.toml]
```

## To-do:
### P1
- [x] Issue builder
- [x] List by category
- [x] Improve single post layout
   - [x] Comments
      - [ ] Captcha layout improvements
   - [x] Image captions
   - [x] Fix to content title
- [x] First archive page
   - [x] Probably the archives are just the lists of posts by category/type
- [ ] Search box
   - [ ] Search results
- [ ] Community page
   - [ ] Pokémon Community Archive
   - [ ] Guestbook
   - [ ] Project ROAR
   - [ ] Product Archive (Here?)
### P2
- [ ] Desktop mode tab bar (dropdown menus?)
- [ ] Log traffic
- [ ] Promoted features
