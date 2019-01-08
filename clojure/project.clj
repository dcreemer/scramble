(defproject scramble 0.1.0-SNAPSHOT"
  :description "Scramble Square Solver, in Clojure"
  :url "http://github.com/dcreemer/scramble-puzzle"
  :license {:name "Eclipse Public License"
            :url "http://www.eclipse.org/legal/epl-v10.html"}
  :dependencies [[org.clojure/clojure "1.9.0"]
                 [org.clojure/math.combinatorics "0.1.4"]
                 [org.clojure/core.async "0.4.490"]]
  :main ^:skip-aot penguin.core
  :target-path "target/%s"
  :profiles {:uberjar {:aot :all}})
