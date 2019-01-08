(ns scramble.core
  (:require [clojure.math.combinatorics :as combo]
            [clojure.core.async :as async])
  (:gen-class))


(def four-exp-9 (apply * (take 9 (repeat 4))))

(defn base4
  "convert a base10 number to base4, represented as a sequence of digits"
  [n]
  (loop [digits [] n n]
    (if (pos? n)
      (recur (conj digits (mod n 4))
             (quot n 4))
      digits)))

(defn nine-digit
  "pad a list of digits with zeros to be 9 digits long"
  [xs]
  (let [need (- 9 (count xs))]
    (concat xs (repeat need 0))))

;; generate every possible combination of 4-way rotations for 9 cards
(def rotations
  (map (comp nine-digit base4)
       (range four-exp-9)))
  

(def matches {:emperorTop      :emperorBottom
              :emperorBottom   :emperorTop
              :adelieTop       :adelieBottom
              :adelieBottom    :adelieTop
              :gentooTop       :gentooBottom
              :gentooBottom    :gentooTop
              :chinstrapTop    :chinstrapBottom
              :chinstrapBottom :chinstrapTop})

;; N, E, W, S, neighbor indexes
(def neighbors {0 [-1  1  3  -1],  1 [-1  2  4  0],  2 [-1 -1  5  1]
                3 [ 0  4  6  -1],  4 [ 1  5  7  3],  5 [ 2 -1  8  4]
                6 [ 3  7 -1  -1],  7 [ 4  6 -1  6],  8 [ 5 -1 -1  7]})

(def opposite {0 2, 1 3, 2 0, 3 1})

(defn get-card-part [deck card-idx part-idx]
  (let [card (nth (:cards deck) card-idx)
        rot  (nth (:rotations deck) card-idx)
        idx  (mod (+ part-idx rot) 4)]
    (nth card idx)))

(defn check-neighbords [deck i j dir]
  ;; check a card and it's neighbor in the given direction
  (let [part (get-card-part deck i dir)
        opart (get-card-part deck j (opposite dir))]
    (= part (matches opart))))

(defn check-card [deck n]
  ;; check all of the neighbors of a single card
  (let [[north east south west] (neighbors n)]
    (and
      (if (< -1 north) (check-neighbords deck n north 0) true)
      (if (< -1 east)  (check-neighbords deck n east 1) true)
      (if (< -1 south) (check-neighbords deck n south 2) true)
      (if (< -1 west)  (check-neighbords deck n west  3) true))))

(defn check-board [board]
  ;; check each card in the board
  (reduce #(and %1 (check-card board %2) true) (range 9)))
                
;; these are the nine cards in the deck
(def cards [[:emperorTop, :adelieBottom, :chinstrapTop, :chinstrapTop]
            [:emperorTop, :gentooTop, :chinstrapTop, :adelieBottom]
            [:emperorTop, :gentooTop, :chinstrapBottom, :adelieTop]
            [:emperorTop, :adelieTop, :gentooBottom, :chinstrapBottom]
            [:emperorBottom, :adelieTop, :chinstrapBottom, :gentooTop]
            [:emperorBottom, :adelieTop, :chinstrapBottom, :gentooBottom]
            [:emperorBottom, :gentooBottom, :chinstrapTop, :adelieBottom]
            [:emperorBottom, :gentooBottom, :chinstrapTop, :adelieBottom]
            [:emperorBottom, :adelieBottom, :gentooTop, :gentooBottom]])

;; some boards
(def lose {:cards cards :rotations [0 0 0 0 0 0 0 0 0]})

(def win  {:cards
           [[:emperorTop, :gentooTop, :chinstrapTop, :adelieBottom]
            [:emperorTop, :adelieTop, :gentooBottom, :chinstrapBottom]
            [:emperorTop, :adelieBottom, :chinstrapTop, :chinstrapTop]
            [:emperorBottom, :adelieTop, :chinstrapBottom, :gentooBottom]
            [:emperorBottom, :gentooBottom, :chinstrapTop, :adelieBottom]
            [:emperorTop, :gentooTop, :chinstrapBottom, :adelieTop]
            [:emperorBottom, :gentooBottom, :chinstrapTop, :adelieBottom]
            [:emperorBottom, :adelieTop, :chinstrapBottom, :gentooTop]
            [:emperorBottom, :adelieBottom, :gentooTop, :gentooBottom]]            
           :rotations [2 2 0 0 0 2 2 2 0]})

;; every possible ordering of all nine cards:
(def orderings (combo/permutations cards))

(def in-chan (async/chan))
(def out-chan (async/chan))

(defn check-all-rotations [cards]
  (loop [rots rotations]
    (let [r  (first rots)
          rr (rest rots)
          board  {:cards cards :rotations r}]
      (if (check-board board)
        (println board))
      (if (not (empty? rr))
        (recur rr)))))

(defn run-workers [num]
  (dotimes [_ num]
    (async/thread
      (loop [board (async/<!! in-chan)
             i 0]
        (check-all-rotations board)
        (if (zero? (mod i 100))
          (println i))
        (recur (async/<!! in-chan) (inc i))))))

(defn gen-work []
  ;; fill the in-chan with work
  (loop [boards orderings]
    (async/>!! in-chan (first boards))
    (recur (rest boards))))

(defn solve []
  (run-workers 2)
  (gen-work))

(defn -main
  [& args]
  (solve))
