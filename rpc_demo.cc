#include <iostream>
#include <string>
#include <functional>
#include <variant>
#include <tuple>
#include <vector>
#include <unordered_map>
#include <stdexcept>

using namespace std;

// Define a variant type to hold any type
using VariantType = std::variant<int32_t, uint32_t, int64_t, uint64_t, float, double, char*, const char*, std::string>;

template <typename... Args>
class rpc_server
{
public:
    rpc_server() = default;
    ~rpc_server() = default;

    void Register(const std::string& name, std::function<void(Args...)>&& callback) {
        mapFuncInfo_[name] = std::move(callback);
    }

    void Call(const std::string& name, const vector<VariantType>& paramList) {
        auto it = mapFuncInfo_.find(name);
        if (it != mapFuncInfo_.end()) {
            auto tup = VariantToTuple<Args...>(paramList);
            std::apply(it->second, tup);
        }else {
            std::cerr << "Function not found: " << name << std::endl;
        }
    }

    template <typename... U>
    auto VariantToTuple(const std::vector<VariantType>& paramList) {
        if (paramList.size() != sizeof...(U)) {
            throw std::invalid_argument("Parameter count mismatch.");
        }
        return ConvertVariantToArgs<std::tuple<U...>>(paramList, std::index_sequence_for<U...>{});
    }

    template <typename Tuple, std::size_t... Index>
    static Tuple ConvertVariantToArgs(const std::vector<VariantType>& paramList, std::index_sequence<Index...>) {
        return std::make_tuple( std::get<std::tuple_element_t<Index, Tuple>>(paramList[Index])...);
    }

private:
    std::unordered_map<std::string, std::function<void(Args...)>> mapFuncInfo_;
};

int main()
{
    std::string func_name = "Hello";
    rpc_server<int, std::string> rpcsvr;
    rpcsvr.Register(func_name, [](int x, const std::string& str)
    {
        std::cout << "Hello call! Integer: " << x << ", String: " << str << std::endl;
    });

    rpcsvr.Call(func_name, {9527, std::string("hello world!")});

    getchar();
    return 0;
}
