Rails.application.routes.draw do
  devise_for :users, controllers: {sessions: "users/sessions"}
  devise_for :installs
  root to: "displays#index"

  resources :vehicles
end