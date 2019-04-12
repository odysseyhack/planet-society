//
//  TransactionOverviewViewController.swift
//  Permission Hub
//
//  Created by Corné on 12/04/2019.
//  Copyright © 2019 Planet. All rights reserved.
//

import UIKit

enum TransactionVerificationStatus {
    case verified, unverified

    var color: UIColor {
        switch self {
        case .verified:
            return PHColors.topaz
        case .unverified:
            return PHColors.red
        }
    }
}

enum TransactionOverviewViewCellType {
    case notification(
        type: TransactionNotificationType,
        text: String)
    case description(text: String)
    case transactionItem(item: TransactionItem)
}

final class TransactionOverviewViewController: UIViewController {

    // MARK: - Private properties

    private lazy var tableView: UITableView = {

        let tableView = UITableView()
        tableView.translatesAutoresizingMaskIntoConstraints = false

        tableView.rowHeight = UITableView.automaticDimension
        tableView.estimatedRowHeight = 66
        tableView.separatorStyle = .none

        tableView.dataSource = self

        tableView.register(
            TransactionTableViewCell.self,
            forCellReuseIdentifier: String(describing: TransactionTableViewCell.self))

        tableView.register(
            TransactionDescriptionTableViewCell.self,
            forCellReuseIdentifier: String(describing: TransactionDescriptionTableViewCell.self))

        tableView.register(
            TransactionNotificationTableViewCell.self,
            forCellReuseIdentifier: String(describing: TransactionNotificationTableViewCell.self))

        return tableView
    }()

    private lazy var bottomStackView: UIStackView = {

        let stackView = UIStackView()
        stackView.translatesAutoresizingMaskIntoConstraints = false

        stackView.axis = .horizontal
        stackView.alignment = .center
        stackView.distribution = .equalCentering
        stackView.spacing = 20

        return stackView
    }()

    private let transaction: TransactionNotification

    private var cells: [TransactionOverviewViewCellType] {

        var cells = [TransactionOverviewViewCellType]()

        cells.append(.notification(
            type: .notification,
            text: "Verified"))
        cells.append(.notification(
            type: .warning,
            text: "Permission warning!"))
        cells.append(.description(text: transaction.reason))

        transaction.items.forEach {
            cells.append(.transactionItem(item: $0))
        }

        return cells
    }

    // MARK: - Initialization

    init(transaction: TransactionNotification) {
        self.transaction = transaction

        super.init(nibName: nil, bundle: nil)
    }

    required init?(coder aDecoder: NSCoder) {
        fatalError("init(coder:) has not been implemented")
    }

    // MARK: - Life cycle

    override func viewDidLoad() {
        super.viewDidLoad()

        view.backgroundColor = PHColors.lightGray

        view.addSubview(tableView)
        view.addSubview(bottomStackView)

        tableView.topAnchor.constraint(equalTo: view.topAnchor).isActive = true
        tableView.leftAnchor.constraint(equalTo: view.leftAnchor).isActive = true
        tableView.rightAnchor.constraint(equalTo: view.rightAnchor).isActive = true
        bottomStackView.heightAnchor.constraint(equalToConstant: 75).isActive = true
        bottomStackView.topAnchor.constraint(equalTo: tableView.bottomAnchor).isActive = true
        bottomStackView.centerXAnchor.constraint(equalTo: view.centerXAnchor).isActive = true
        bottomStackView.bottomAnchor.constraint(equalTo: view.bottomAnchor).isActive = true

        for i in 0..<2 {

            let button = UIButton()
            button.translatesAutoresizingMaskIntoConstraints = false

            button.widthAnchor.constraint(equalToConstant: 92).isActive = true
            button.heightAnchor.constraint(equalToConstant: 34).isActive = true

            button.titleLabel?.font = PHFonts.regular()

            let color = i == 0 ? PHColors.greyishBrown : PHColors.topaz
            button.setTitleColor(color, for: .normal)
            button.backgroundColor = .white

            button.layer.cornerRadius = 5
            button.layer.borderWidth = 1
            button.layer.borderColor = color.cgColor

            let title = i == 0 ? "Decline" : "Continue"
            button.setTitle(title, for: .normal)
            
            bottomStackView.addArrangedSubview(button)
        }
    }
}

extension TransactionOverviewViewController: UITableViewDataSource {

    func tableView(_ tableView: UITableView, numberOfRowsInSection section: Int) -> Int {
        return cells.count
    }

    func tableView(_ tableView: UITableView, cellForRowAt indexPath: IndexPath) -> UITableViewCell {

        switch cells[indexPath.row] {
        case .notification(let type, let text):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionNotificationTableViewCell.self),
                for: indexPath) as! TransactionNotificationTableViewCell

            cell.configure(withType: type, andText: text)

            return cell

        case .description(let text):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionDescriptionTableViewCell.self),
                for: indexPath) as! TransactionDescriptionTableViewCell

            cell.configure(withText: text)

            return cell

        case .transactionItem(let item):

            let cell = tableView.dequeueReusableCell(
                withIdentifier: String(describing: TransactionTableViewCell.self),
                for: indexPath) as! TransactionTableViewCell

            let viewModel = TransactionTableViewCellViewModel(
                image: UIImage(),
                title: item.item,
                subtitle: (item.fields as NSArray).componentsJoined(by: ", "))
            cell.configure(withViewModel: viewModel)

            return cell
        }
    }
}

extension TransactionOverviewViewController: UITableViewDelegate {

    func tableView(_ tableView: UITableView, didSelectRowAt indexPath: IndexPath) {

        switch cells[indexPath.row] {
        case .notification(let type, let text):

            let viewController = UITableViewController()
            navigationController?.pushViewController(viewController, animated: true)

        default:
            break
        }
    }
}
