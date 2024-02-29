#include <sstream>
#include <map>
#include <vector>
using namespace std;

class Solution {
public:
    bool wordPattern(string p, string s) {
        istringstream iss(s);
        string word;
        vector<string> words;

        while(iss >> word){
            words.push_back(word);
        }

        if (p.size() != words.size()) {
            return false;
        }

        map<char, string> p_to_s;
        map<string, char> s_to_p;

        for (int i = 0; i < p.size(); i++){

            if (p_to_s.find(p[i]) != p_to_s.end()){//パターン文字がマップに見つかった
                if (p_to_s[p[i]] != words[i]){//それが対応するものじゃなかった
                    return false;
                }
            }
            if (s_to_p.find(words[i]) != s_to_p.end()) {//単語がマップに見つかった
                if (s_to_p[words[i]] != p[i]) {//それが対応するものじゃなかった
                    return false;
                }
            }

            // どっちもマップから見つからなかったとき、新しいパターン文字や単語をマップに追加
            p_to_s[p[i]] = words[i];
            s_to_p[words[i]] = p[i];
        } //繰り返す

        return true;

    }
};