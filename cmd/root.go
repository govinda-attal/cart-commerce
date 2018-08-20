package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/urfave/negroni"

	"github.com/govinda-attal/cart-commerce/internal/cartcom"
	"github.com/govinda-attal/cart-commerce/internal/promo"
	"github.com/govinda-attal/cart-commerce/internal/provider"
	"github.com/govinda-attal/cart-commerce/pkg/core/httputil"
)

var (
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cartcom",
	Short: "Starts microservice",
	Run:   startServer,
}

func startServer(cmd *cobra.Command, args []string) {
	provider.Setup()

	db := provider.DB()

	promoApi := promo.NewApiImpl()
	cartApi := cartcom.NewApiImpl(db, promoApi)

	cartHandler := cartcom.NewHandler(cartApi)
	ruleHandler := promo.NewHandler(promoApi)

	r := mux.NewRouter()

	r.PathPrefix("/1.0/api/").Handler(http.StripPrefix("/1.0/api", http.FileServer(http.Dir("./api"))))
	//r.PathPrefix("/1.0/rules/").Handler(http.StripPrefix("/rules", http.FileServer(http.Dir("./rules"))))

	r.HandleFunc("/1.0/rules", httputil.WrapperHandler(ruleHandler.Fetch)).Methods("GET")

	r.HandleFunc("/1.0/ecart/{cartId}", httputil.WrapperHandler(cartHandler.FetchCartItems)).Methods("GET")
	r.HandleFunc("/1.0/ecart", httputil.WrapperHandler(cartHandler.UpdateCartItem)).Methods("POST")
	r.HandleFunc("/1.0/ecart/{cartId}", httputil.WrapperHandler(cartHandler.UpdateCartItem)).Methods("PUT")
	r.HandleFunc("/1.0/ecart/{cartId}", httputil.WrapperHandler(cartHandler.UpdateCartState)).Methods("PATCH")

	h := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH"},
	}).Handler(r)
	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.UseHandler(h)

	srv := &http.Server{
		Addr:         "0.0.0.0:9080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      n,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	provider.Cleanup()
	srv.Shutdown(ctx)
	log.Println("cartcom server shutdown ...")
	os.Exit(0)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config/app-config.yml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if len(cfgFile) != 0 {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigFile("./config/app-config.yml")
	}
	viper.AutomaticEnv()
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}
